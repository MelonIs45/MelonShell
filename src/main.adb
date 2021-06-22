with Ada.Text_IO;
use Ada.Text_IO;

procedure Main is
begin

   while (True) loop

      Put("$> ");
      declare
         input : String := Get_Line;
      begin
         -- once you get input, carry on here

         Put_Line (input);
      end;

   end loop;
   null;
end Main;
